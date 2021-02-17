# SilverBroccoli

A small pet project to dive into Golang language. The name was taken from random name that gives Github when create a new repository, sorry for my lack of creativity.

![](https://github.com/sjdiego/silverbroccoli/workflows/Go/badge.svg)

## What it does
It stores posts from Reddit into MongoDB collections. Then the posts can be posted to Telegram channels.

## Requirements
It needs ffmpeg to merge video and audio files from Reddit as they send it separately.

## How to use
1. Copy `config/config.example.yml` to `config/config.yml` and change your MongoDB data and Telegram API ket.

2. Build the code and execute with any of this flags:

   - __-addchan__: prompts for data to store a Telegram channel.
   - __-addsub__: prompts for data to store a Subreddit to save its posts
   - __-fetch__: fetch and store posts from stored subreddits.
   - __-publish__: send stored and unpublished posts from MongoDB to stored Telegram channels. Use __-num `<int>`__ to change the number of posts to send.

## Notes
- The user input of prompts is not validated, so if you type incorrect values it will fail.
- If more than one Telegram channel has the same subreddit than another Telegram channel its posts won't be send as the previous channel has already set the posts as sent.
- Lack of tests.
- It only sends posts with videos.
