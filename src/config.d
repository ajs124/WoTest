{- vim: set ft=dhall: -}
let enums = (./defs.d).enums
{- let types = (./defs.d).types -}
let implemenntationsDir = "../implementations"
in {
    logLevel = enums.logLevels.debug,
    testResults = "../results.json",
    testsDir = "../tests",
    implementationsDir = implemenntationsDir,
    implementations = [
        {
            name = "node-wot",
            path = implemenntationsDir ++ "/node-wot-install/node_modules",
            runtime = enums.runtimes.node
        },
        {
            name = "wot-py",
            path = implemenntationsDir ++ "/wot-py/install",
            runtime = enums.runtimes.python
        },
        {
            name = "sane-city",
            path = "", -- TODO
            runtime = enums.runtimes.java
        }
    ]
}
