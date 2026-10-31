# PROGRESS

One row per finished pattern. Add yours after step 6 of [HOWTO.md](HOWTO.md).
"What I got wrong" is the useful column, so do not leave it blank.

Next: **A1 (SRP)**


| ID | Date  | What I got wrong                                                                                                                | My one-line "don't" |
| -- | ----- | ------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| 1  | 26/10 | Permission of as os.ModeAppend but it should be 0o644 and it did't use %w for error wrap                                        |                     |
| 2  | 27/10 | Didn;t added lock while writing dict, name chose were not good, and didn't added the condition to avoid duplicate tier addition |                     |
| 3  | 27/10 | Missed to check sentinal error, wrote not a clean code, like error message can be more informative                              |                     |
| 4  | 28/10 | Used %s in place %q, error wrap will have not needed info as I have added "error while getting user %w, err"                              |                     |
| 5  | 29/10 | Used strings pkg to trim, its like leak provider should clean and provide output, consumer should use it as it is.                              |                     |
| 6  | 29/10 | Used %w for Get func since its the deepest layer hence it not be used there                              |                     |
