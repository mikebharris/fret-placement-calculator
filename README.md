# Fret Placement Calculator AWS Lambda service

This simple AWS Lambda function depends on my music theory Go module at https://github.com/mikebharris/music.  All the computation is done in that module; this project is just a thin wrapper to expose it as a web service.

The service outputs where to place the frets on a fretboard of a stringed instrument, effectively where to stop the strings, for
various tunings, including:

* Just Intonation by various means
* Ptolemy's Intense Diatonic scale in various modes
* Quarter-Comma Meantone
* Extended Quarter-Comma Meantone
* Any Equal Temperament (12, 19, 23, 31, 53, 55, etc)
* Pythagorean
* Turkish Saz
* Bach's Well Temperament (as decoded by Bradley Lehman)

The fret positions are agnostic of the actual open string tuning, tension of the string, type of instrument, etc.

## Examples

<details>
 <summary><code>GET</code> <code><b>/scaleLength={scaleLength}</b></code> <code>(returns fret positions for just intonation and scale length of {scaleLength}</code></summary>

##### Parameters

> | name           | type     | data type | default | description                                                                                                 |
> |----------------|----------|-----------|---------|-------------------------------------------------------------------------------------------------------------|
> | `scaleLength`  | required | float64   |         | The scale length from nut to bridge (saddle)                                                                |
> | `tuningSystem` | required | string    |         | Tuning to use (just, meantone, pythagorean, equal, ptolemy, saz).  Defaults to a chromatic Just tuning.     |
> | `diatonicMode` | optional | string    | Ionian  | Produce a diatonic scale instead of chromatic in the specified musical mode (ionian, dorin, phryggian, etc) |
> | `limit`        | optional | int       | 5       | Limit for just intonation (prime number, such as 3, 5, 11, etc_ - tuningSystem = 'justFromRatios'           |
> | `division`     | optional | int       | 31      | Number of divisions of the octave for equal temperament                                                     |
> | `octaves`      | optional | int       | 1       | Number of octaves of frets to compute                                                                       |

##### Values for `tuningSystem`

> | value                       | description                                                                         |
> |-----------------------------|-------------------------------------------------------------------------------------|
> | `justFromRatios`            | 5-limit Just Intonation derived from pure ratios                                    |
> | `just5limitFromPythagorean` | 5-limit Just Intonation derived from tweaking Pythagorean scale by a syntonic comma |
> | `meantone`                  | Quarter-Comma Meantone                                                              |
> | `extendedMeantone`          | Extended Quarter-Comma Meantone                                                     |
> | `bachWellTemperament`       | Bach's Well Temperament (as decoded by Bradley Lehman)                              |
> | `pythagorean`               | Pythagorean 3-limit just tuning                                                     |
> | `equal`                     | Equal Temperament                                                                   |
> | `ptolemy`                   | Ptolemy's Intense Diatonic tuning                                                   |
> | `saz`                       | Turkish Saz tuning                                                                  |

##### Responses

> | http code | content-type       | response                                 |
> |-----------|--------------------|------------------------------------------|
> | `200`     | `application/json` | JSON object                              |
> | `422`     | `application/json` | `{"code":"422","message":"Bad Request"}` |

##### Example cURL

Compute Ptolemy's Intense Diatonic tuning for a scale length of 570mm:

> ```shell
>  curl -X GET -H "Content-Type: application/json" https://someawsgeneratedlambdaid.lambda-url.us-east-1.on.aws/?scaleLength=570&tuningSystem=ptolemy
> ```

````json
{
  "system": "Ptolemy Intense Diatonic",
  "description": "Fret positions based on Ptolemy's 5-limit intense diatonic scale in Ionian mode.",
  "scaleLength": 570,
  "frets": [
    {
      "label": "9:8",
      "position": 63.33,
      "comment": "Pythagorean (Greater) Major Second",
      "interval": "9:8"
    },
    {
      "label": "5:4",
      "position": 114,
      "comment": "Major Third",
      "interval": "10:9"
    },
    {
      "label": "4:3",
      "position": 142.5,
      "comment": "Perfect Fourth",
      "interval": "16:15"
    },
    {
      "label": "3:2",
      "position": 190,
      "comment": "Perfect Fifth",
      "interval": "9:8"
    },
    {
      "label": "5:3",
      "position": 228,
      "comment": "Major Sixth",
      "interval": "10:9"
    },
    {
      "label": "15:8",
      "position": 266,
      "comment": "Just Major Seventh",
      "interval": "9:8"
    },
    {
      "label": "2:1",
      "position": 285,
      "comment": "Perfect Octave",
      "interval": "16:15"
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
AWS_ACCESS_KEY_ID=???? AWS_SECRET_ACCESS_KEY=???? go run pipeline.go --account-number=123456789012 --app-name=fret-placement-calculator --environment=prod --region=us-east-1 --stage=apply --confirm=true
````

Refer to the documentation in that project for more details on how to use the deployment helper tool.
