API:
- `/api/test/add`: Add a test with format:

```
{
    "name": "Among us quiz" // Name of the test
    "start": 1668006000 // Start time, in unix time
    "end": 1668009600 // End time, in unix time
    "questions": [
        {
            "content": "Is red sus?" // Content of the question
            "choices": [
                {
                    "content": "Yes" // Content of the choice
                },
                {
                    "content": "No" 
                }
            ]
        },
    ]
}
```
