# Little LATAM job applicant information compilation using Revelo's Linkedin job postings

This study consisted in scraping all of the company's (Revelo's) job postings.

The company seems to have hundreds of job postings and regularly takes the job posting down
only to resubmit it as a new job posting. This seems to be done programatically 
as all programmer job postings have the same age.

## Introduction 
Why they do this is of not interest to us. What is of interest is the valuable
information in the number of applicants in each job posting. This gives us
information on how many developers know or are willing to work using a language
given a certain salary. We have access to the following information for each job posting:

- The technology desired for the job posting
- The location **\***
- The number of applicants
- The salary range

**\*** The job postings in themselves are labelled as remote so the location will
not reliably tell us where applicants live. That said, it does seem like the Linkedin advertisement for the job posting
tends to be shown to people in the posting location.

## Methodology
Download HTML directly from linkedin for all revelo available jobs.

Process and generate reports using program [`main.go`](./main.go).

## Results

- See [Report by job description](./report_by_job.md)
- See [Report by location](./report_by_location.md)
- Raw results in [`results.csv`](./results.csv) file


