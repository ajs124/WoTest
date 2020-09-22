{- vim: set ft=dhall: -}
{
    enums = {
        logLevels = {
            panic =  5,
            fatal =  4,
            error =  3,
            warn  =  2,
            info  =  1,
            debug =  0,
            trace = -1
        },
        runtimes = {
            node = 0,
            python = 1,
            java = 2
        },
        tests = {
            first = 0
        }
    }
    {- types = {
        test : Type = {
            path: Text,
            type: Number,
            timeoutSec: Number
        }
    } -}
}
