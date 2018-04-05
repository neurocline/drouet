# Hugo conversion log

## Methodology

1. Fork `gohugoio/hugo` ==> `neurocline/drouet`.
2. the branch `hugo:drouet-translation` has code removed from it as is it added to `drouet`. This is needed because Drouet is not a line-for-line copy of Hugo.
3. `drouet-translation` is periodically rebased on top of `hugo:master`; conflicts show us changes in code that was already moved over.
