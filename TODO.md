# TODO list

* Use [govaluate] or similar to allow conditional enable/disable based on request parameters. This would require an additional argument to `feature.Get` (or a second function) to pass additional parameters, and an expansion of the Feature type to allow for things more complex than simply "100% on" or "100% off".
* Percentage-based feature flags.

[govaluate]: https://github.com/Knetic/govaluate