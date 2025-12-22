# Fret Placement Calculator AWS Lambda service

Outputs where to place the frets on a fretboard of a stringed instrument, effectively where to stop the strings, for
various tunings, including:

* Just Intonation
* Quarter-Comma Meantone
* Extended Quarter-Comma Meantone
* Any Equal Temperament (12, 19, 23, 31, 53, 55, etc)
* Pythagorean
* Turkish Saz

With the exception of the quarter-comma meantone calculator, the rest are agnostic of the actual open string tuning,
tension of the string, type of instrument, etc.

## Examples

<details>
 <summary><code>GET</code> <code><b>/scaleLength={scaleLength}</b></code> <code>(returns fret positions for just intonation and scale length of {scaleLength}</code></summary>

##### Parameters

> | name                    | type     | data type | default    | description                                                                                                 |
> |-------------------------|----------|-----------|------------|-------------------------------------------------------------------------------------------------------------|
> | `scaleLength`           | required | float64   |            | The scale length from nut to bridge (saddle)                                                                |
> | `tuningSystem`          | optional | string    | just       | Tuning to use (just, meantone, pythagorean, equal, ptolemy, saz).  Defaults to a chromatic Just tuning.     |
> | `diatonicMode`          | optional | string    | Ionian     | Produce a diatonic scale instead of chromatic in the specified musical mode (ionian, dorin, phryggian, etc) |
> | `justSymmetry`          | optional | string    | asymmetric | Type of major seconds and minor sevenths to use in just scale                                               |
> | `extendMeantone`        | optional | bool      | false      | Extend the meantone scale                                                                                   |
> | `octaveDivisions`       | optional | int       | 31         | Number of divisions of the octave for equal temperament                                                     |
> | `octaves`               | optional | int       | 1          | Number of octaves of frets to compute                                                                       |

##### Values for `tuningSystem`

> | value                       | description                                                                         |
> |-----------------------------|-------------------------------------------------------------------------------------|
> | `just5limitFromRatios`      | 5-limit Just Intonation derived from pure ratios                                    |
> | `just5limitFromPythagorean` | 5-limit Just Intonation derived from tweaking Pythagorean scale by a syntonic comma |
> | `meantone`                  | Quarter-Comma Meantone                                                              |
> | `pythagorean`               | Pythagorean 3-limit just tuning                                                     |
> | `equal`                     | Equal Temperament                                                                   |
> | `ptolemy`                   | Ptolemy's Intense Diatonic tuning                                                   |
> | `saz`                       | Turkish Saz tuning                                                                  |

##### Values for `justSymmetry`

> | value        | description                                                               |
> |--------------|---------------------------------------------------------------------------|
> | `asymmetric` | Use asymmetric scale with greater major seconds and lesser minor sevenths |
> | `symmetric1` | Use symmetric scale with lesser major seconds and greater minor sevenths  |
> | `symmetric2` | Use symmetric scale with greater major seconds and lesser minor sevenths  |


##### Responses

> | http code | content-type       | response                                 |
> |-----------|--------------------|------------------------------------------|
> | `200`     | `application/json` | JSON object                              |
> | `422`     | `application/json` | `{"code":"422","message":"Bad Request"}` |

##### Example cURL

Compute Ptolemy's Intense Diatonic tuning for a scale length of 546mm:

> ```shell
>  curl -X GET -H "Content-Type: application/json" https://someawsgeneratedlambdaid.lambda-url.us-east-1.on.aws/?scaleLength=546
> ```

````json
{
  "system": "5-limit Just Intonation",
  "description": "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from applying syntonic comma to Pythagorean ratios.",
  "scaleLength": 546,
  "frets": [
    {
      "label": "9:8",
      "position": 60.67
    },
    {
      "label": "5:4",
      "position": 109.2
    },
    {
      "label": "4:3",
      "position": 136.5
    },
    {
      "label": "3:2",
      "position": 182
    },
    {
      "label": "5:3",
      "position": 218.4
    },
    {
      "label": "16:9",
      "position": 238.88
    },
    {
      "label": "15:8",
      "position": 254.8
    },
    {
      "label": "2:1",
      "position": 273
    }
  ]
}
````

</details>

## Building and provisioning

To build this project, copy the
tool https://github.com/mikebharris/aws-deployment-pipeline-orchestration-helper-tool/blob/main/pipeline.go the project
at https://github.com/mikebharris/aws-deployment-pipeline-orchestration-helper-tool into the top-level directory, and
do:

```shell
go mod tidy
go run pipeline.go --help
```

An example Terraform build and deploy command line:

```shell
go run pipeline.go --stage=build
AWS_ACCESS_KEY_ID=???? AWS_SECRET_ACCESS_KEY=???? go run pipeline.go --account-number=123456789012 --app-name=fret-placement-calculator --environment=prod --region=us-east-1 --stage=plan
````

Refer to the documentation in that project for more details on how to use the deployment helper tool.
