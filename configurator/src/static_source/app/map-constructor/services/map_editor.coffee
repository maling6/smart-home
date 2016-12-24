angular
.module('angular-map')
.factory 'mapEditor', ['$rootScope', '$compile', 'mapFullscreen', 'mapPanning', '$templateCache', 'mapLayer', 'mapElement'
  ($rootScope, $compile, mapFullscreen, mapPanning, $templateCache, mapLayer, mapElement) ->
    class mapEditor

      constructor: ()->
        @scope.addLayer = @addLayer
        @scope.removeLayer = @removeLayer
        @scope.selectLayer = @selectLayer
        @scope.selectElement = @selectElement
        @scope.addElement = @addElement
        @scope.current_layer = null
        @scope.sortLayers = @sortLayers
        @scope.sortElements = @sortElements

      loadEditor: (c)=>
        # container
        # --------------------
        _container = angular.element(c)
        template = $templateCache.get('/map-constructor/templates/map_editor.html')
        angular.element(_container).append($compile(template)(@scope))

        @container = _container.find(".map-editor")
        @wrapper = _container.find(".map-wrapper")

        # fullscreen
        # --------------------
        @fullscreen = new mapFullscreen(@wrapper, @scope)

        # resizable
        # --------------------
        if @wrapper.resizable('instance')
          @wrapper.resizable('destroy')
        @wrapper.resizable
          minHeight: @scope.settings.minHeight
          minWidth: @scope.settings.minWidth
          grid: @scope.settings.grid
          handles: 's'

        @panning = new mapPanning(@container, @scope, @wrapper)
        @wrapper.find(".page-loader").fadeOut("fast")

        return

      addLayer: ()=>
        if !@model?.layers
          @model.layers = []

        layer = new mapLayer(@scope)
        layer.map_id = @id
        @model.layers.push layer

      removeLayer: (_layer)=>
        index = @model.layers.indexOf(_layer)
        if index > -1
          @model.layers.splice(index, 1)
          @scope.current_layer = null
        return

      selectLayer: (layer, $index)=>
        if @scope.current_layer == layer
          @scope.current_layer = null
        else
          @scope.current_layer = layer

      selectElement: (element, $index)=>
        if @scope.current_element == element
          @scope.current_element = null
        else
          @scope.current_element = element

      addElement: ()=>
        return if !@scope.current_layer
        @scope.current_layer.addElement new mapElement(@scope)

      removeElement: (_element)=>
        index = @scope.current_element.elements.indexOf(_element)
        if index > -1
          @scope.current_element.elements.splice(index, 1)
          @scope.current_element = null
        return

      sortLayers: ()=>
        weight = 0
        for layer in @model.layers
          layer.weight = weight
          weight++

      sortElements: ()=>
        return if !@scope.current_layer
        weight = 0
        E
        for element in @scope.current_layer.elements
          element.weight = weight
          weight++


    mapEditor
]