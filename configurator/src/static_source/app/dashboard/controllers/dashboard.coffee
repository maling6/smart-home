angular
.module('appControllers')
.controller 'dashboardCtrl', ['$scope', 'Notify'
($scope, Notify) ->
  vm = this

  vm.menu =
    'name': 'Main Menu'
    'items': [
      {
        'label': 'Dashboard'
        'route': 'dashboard.index'
        'icon': 'fa fa-home'
      }
      {
        'label': 'Nodes'
        'route': 'dashboard.node'
        'icon': 'fa fa-sitemap'
        'items': [
          {
            'label': 'All nodes'
            'route': 'dashboard.node'
          }
          {
            'label': 'Add new'
            'route': 'dashboard.node_new'
          }
        ]
      }
    ]

  vm
]