let GitHubActions = (../imports.dhall).GitHubActions

let Setup = ../setup.dhall

in  Setup.MakeJob
      Setup.JobArgs::{
      , name = "c-format"
      , additionalSteps =
        [ GitHubActions.Step::{
          , name = Some "c-format"
          , uses = Some "jidicula/clang-format-action@v4.0.0"
          , `with` = Some
              ( toMap
                  { source = "."
                  , clangFormatVersion = "11"
                  , fallback-style = "Mozilla"
                  }
              )
          }
        ]
      }
