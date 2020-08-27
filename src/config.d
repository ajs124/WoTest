{- vim: set ft=dhall: -}
let
    logLevels = {
        panic =  5,
        fatal =  4,
        error =  3,
        warn  =  2,
        info  =  1,
        debug =  0,
        trace = -1
    }
let runtimes = {
        node = 0,
        python = 1,
        java = 2
    }
let ImplementationsDir = "../implementations"
in {
    LogLevel = logLevels.debug,
    TestsDir = "../tests",
    ImplementationsDir = ImplementationsDir,
    Implementations = [
        {
            name = "node-wot",
            path = ImplementationsDir ++ "/node-wot-install/node_modules",
            runtime = runtimes.node
        },
        {
            name = "wot-py",
            path = ImplementationsDir ++ "/wot-py/install",
            runtime = runtimes.python
        },
        {
            name = "sane-city",
            path = "", -- TODO
            runtime = runtimes.java
        }
    ]
}
