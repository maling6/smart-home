package workflow

const coffeeScript1 = `
main =->
    wf = IC.Workflow()
    print wf.getName() + " script1"
`

const coffeeScript2 = `
main =->
    wf = IC.Workflow()
    print wf.getName() + " script2"
`

const coffeeScript3 = `
main =->
    print "workflow scenario script 3"

on_enter =->
    wf = IC.Workflow()
    print "enter to " + wf.getName() + " " + wf.getScenarioName()
    wf.setVar("var1", 123)
    wf.setScenario("wf_scenario_2")

on_exit =->
    wf = IC.Workflow()

    print "exit from " + wf.getName() + " " + wf.getScenarioName()
    print "description " + wf.getDescription()
`



const coffeeScript4 = `
main =->
    print "workflow scenario script 4"

on_enter =->
    wf = IC.Workflow()

    print "enter to " + wf.getName() + " " + wf.getScenarioName()
    var1 = wf.getVar("var1")
    print "var: " + var1

on_exit =->
    wf = IC.Workflow()

    print "exit from " + wf.getName() + " " + wf.getScenarioName()
    print "description " + wf.getDescription()
`

const coffeeScript5 = `
# test 2

print "run message emitter script (script 5)"
print "message:", message
print "ENV", ENV.a

print "message.ENV:", message.getVar "ENV"

IC.store ENV.a
`

const coffeeScript6 = `
# test 2

ENV = {"a": "b"}
print "run message handler script (script 6)"

message.setVar("ENV", ENV) 
`

const coffeeScript7 = `
# test 2
# workflow script

print "run workflow script (script 7)"
`

const coffeeScript8 = `
print "run workflow script (script 8)"
`

const coffeeScript9 = `
print "run workflow script (script 9)"
`

const coffeeScript10 = `
print "run workflow script (script 10)"
`

const coffeeScript11 = `
print "run workflow script (script 11)"
`