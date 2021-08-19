let GitHubActions = (../imports.dhall).GitHubActions

let Setup = ../setup.dhall

in  Setup.MakeJob
      Setup.JobArgs::{
      , name = "go-test"
      , additionalSteps =
        [ GitHubActions.Step::{
          , name = Some "go-test"
          , run = Some "ci/go-test.sh"
          }
        ]
      }
