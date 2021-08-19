let GitHubActions = (../imports.dhall).GitHubActions

let Setup = ../setup.dhall

in  Setup.MakeJob
      Setup.JobArgs::{
      , name = "c-format"
      , additionalSteps =
        [ GitHubActions.Step::{
          , name = Some "c-format"
          , uses = Some "DoozyX/clang-format-lint-action@v0.11"
          , `with` = Some (toMap { source = ".", clangFormatVersion = "11" })
          }
        ]
      }
