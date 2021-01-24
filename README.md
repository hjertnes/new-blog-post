# README

This is a simple utility program I wrote to make it faster to make new blog posts for my hugo blog. 

It has two sub commands: post and status. A post contains a tilte while a status don't. You can pass --photo to either and it will add markdown links to all the pictures in your current dir the static of your site.

All commands opens the created file in EDITOR

## Installation
You need to have go installed. And you need to run: `make build` to build it and then copy new-blog-post to somewhere in your PATH

## Config 

You need to set the ENV variable `HUGO_POST_COLLECTION` to where you want your files to end up.