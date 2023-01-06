API:
- `/api/test/add`: Add a test with format:

```json
{
    "name": "Among us quiz",
    "start": 1668006000, // Unix time
    "end": 1668009600, // Unix time
    "questions": [
        {
            "content": "Is red sus?",
            "choices": [
                {
                    "content": "Yes",
                    "is_answer": true,
                },
                {
                    "content": "No",
                    "is_answer": false,
                }
            ]
        },
    ]
}
```

- `/api/test/delete`: Delete a test with given id