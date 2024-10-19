# gotaglint

gotaglint is a program to check whether you have have defined all required build tags.

## Install

```sh
go install github.com/hrvadl/gotaglint/cmd/gotaglint
```

## Usage

To scan all packages run:

```sh
gotaglint --match=".*_test.go" --buildtags="integration,\!integration" ./...
```

Where `--match` is the regex of files to check and `--buildtags` is required build tags, separated by comma.

## Example

Let's say you have following code:

```go
package examples

import "testing"

func Test_sum(t *testing.T) {
	t.Parallel()
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "ShouldSumCorrectly",
			args: args{
				a: 1,
				b: 1,
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sum(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

And run `gotaglint` with the following parameters:

```sh
gotaglint --match=".*_test.go" --buildtags="integration,\!integration" ./...
```

`gotaglint` will produce the following output:

```sh
/Users/vadym.hrashchenko/go/gotaglint/examples/example_test.go:3:9: matching build tag is not found
```

## Rules

It highlights you if you forgot to specify build tags. This can be particularly useful if you separate
unit and integration tests in your `Go` code by build tags. For example:

```sh
//go:build integration

OR

//go:build !integration
```
