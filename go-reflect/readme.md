
# go-reflect

This is an experiment to deal with loading records into a structure.

Ideally we can use the `json` package as a shortcut to achieving a simple mapping.

However, in actual use we run into edge-cases like silently discarding potentially valid values and related casting errors.


## problem

To summarize, the problem is that we often want to read data from multiple sources, and not every source provides valid `json` data.

An actual example would be loading configuration data from three sources:

- configuration file
- command line options
- environment variables

While the configuration file may conform to valid json data types, the environment variables and command line options are always strings.


### `json`

I figured the fastest way to figure this out would be to reverse engineer `json.Unmarshal`; I was wrong but I did learn quite a bit.

The chain of function calls is very complex and depends on the `reflect` package for all of the heavy-lifting, including parsing the json tags used to map properties.

The major take-aways are that json deals with three raw data types:

- numbers (float64)
- booleans
- strings

_The main problem is that the casting logic is limited to expected types; a `float64` can be cast into any other numeric type, but a string cannot.  The same applies to boolean casting._

Thus it only solves for three scenarios:

- number to number
- boolean to boolean
- string to string

The `json` package does map direct property matches in addition to overriding matching property values with values in matching tags.


### edge-cases

The following edge cases are not solved for implicitly:

- string to boolean
- string to number
- number to string
- boolean to string
- boolean to number
- string boolean to number

_The first two are the most important for our aforementioned use case, while the others are simply possible to solve for with minimal extra effort._


## solution

**To summarize, we can still use the `json` package in conjunction with the `reflect` package to pre-emptively cast by matching tags; best of all we only have to work with three data types instead of all go primitives.**

This fix requires the structure you are placing data into, and must be run prior against a `map[string]interface{}` to fix data, marshal the fixed types into json, then unmarshal that onto the structure.

We apply casting rules in two cases:

- no tags apply to the property
- the tag is matched

**We do not cast when a property has a tag, because we may blow up on this scenario:**

	type Convoluted struct {
		Name     string `json:"name"`
		DomainId int    `json:"Name"`
	}

Keeping in-line with the `json` package, our solution silently discards casting errors (_so only when we cannot parse the expected result do we discard it_).

_This implementation is not recommended for applications that require high throughput; the behavior involves multiple conversions which could be sped up by replacing the `json` package entirely and writing the mappings directly alongside the extra casting logic._


## code

The code demonstrates default behavior then edge cases.

The default behavior shows that casting works correctly when the data types align.

Additionally direct property matches will work, even when tags exist.

When a tag value and property match value both exist, the tag match will replace the property match.

The edge cases demonstrate the use of a reflection loop with the structure as the context, and handles all the basic conversion to valid json data types.

Once the conversion occurs we turn the map into json, then the json into the structure.


## usage

For convenience this example includes a full test and benchmark suite:

	go test -v -bench=.


## benchmarks

**These are from my own personal machine and may vary by platform; ideally these serve as a basic idea as to the complexity each step adds to conversion.**

Performance with only the initial property loop only the expected map properties:

	200000       	      6164 ns/op

Performance after adding all unexpected properties to the map:

	100000       	     12161 ns/op

After adding a second loop:

	100000       	     14868 ns/op

After adding conditional check that parses tags and matches property names:

	30000       	     50364 ns/op

After adding all conversions:

	20000       	     59627 ns/op

Benchmarks with added recursion to support depth:

	10000       	    117798 ns/op

Adding early-termination:

	10000       	    118405 ns/op

_The early termination appears to have no impact on performance, so it may make more sense to remove it as it only helps for string to string and bool to bool scenarios due to the number of numeric combinations._

Post-optimization benchmarks:

	50000       	     41016 ns/op

_Moving the reflect calls outside the k/v map loop cut the time by 2/3, and simplified the code_

This time we added a generic isNumeric function and replaced embedded `switch` with single-depth `if` blocks; this solves our missing early termination function and makes the code more succinct:

	30000       	     39788 ns/op


## bugs

There are issues with the `json` package in the latest version of go (currently tested with 1.7) which cause valid json with invalidly mapped data to cause other valid fields to be discarded.

Invalid mappings could be fields that don't exist on the target structure, or unexpected types coming from json.

The code here can resolve the type-matching concerns, but the non-existent fields may still lead to problems.
