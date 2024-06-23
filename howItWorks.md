#### Usage

This script can be run from the command line and accepts several arguments that specify the mode of operation and processing details. Below is the general syntax to run the script:

```bash
go run <script-name>.go <inputLink> [numThreads] [mode]
```

- `<script-name>.go` should be replaced with the actual filename of the Go script.
- `<inputLink>` is a mandatory argument that specifies the input file. Valid options for this are:
    - `xsmall`
    - `small`
    - `medium`
    - `large`

- `[numThreads]` is an optional argument that specifies the number of threads to be used for parallel processing. This should be a positive integer.

- `[mode]` is an optional argument that determines the type of parallel processing:
    - `p` for standard parallel processing
    - `q` for a work-queue based parallel processing

#### Examples

1. **Sequential Processing:**
   If you want to process the data sequentially, you only need to provide the `inputLink`. For example:

   ```bash
   go run script.go xsmall
   ```

2. **Parallel Processing:**
   To process data in parallel, specify the number of threads and the mode. For example:

   ```bash
   go run script.go small 4 p
   ```

   This will process the `small` input file using 4 threads in standard parallel mode.

3. **Work-Queue Based Parallel Processing:**
   If you prefer work-queue based parallel processing, use `q` as the mode. For example:

   ```bash
   go run script.go medium 3 q
   ```

   This command processes the `medium` input file using 3 threads in a work-queue based parallel mode.

#### Error Handling

If there is an issue with the number of threads (e.g., non-integer input), the script will output an error message and terminate:

```
Error converting number of threads: [error details]
```

Ensure that the number of threads is a valid integer.


