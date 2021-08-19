let GitHubActions = (../imports.dhall).GitHubActions

let Setup = ../setup.dhall

in  Setup.MakeJob
      Setup.JobArgs::{
      , name = "golangci-lint"
      , additionalSteps =
        [ GitHubActions.Step::{
          , name = Some "golangci-lint"
          , uses = Some "golangci/golangci-lint-action@v2.5.2"
          , `with` = Some (toMap { version = "v1.39.0" })
          }
        ]
      }
