
# Audio go stream

This web service enables the user to stream the audio of any web page.
To use the service, open the following url in any browser:

```
<host>/audio/<url-to-webpage>
```

Where `<url-to-webpage>` can be any page, such as `https://www.youtube.com/watch?v=dQw4w9WgXcQ`.

The result could look something like:

```
http://localhost:8910/audio/https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

# How it works

the web service opens a browser (chromium) and uses PulseAudio (https://www.freedesktop.org/wiki/Software/PulseAudio/)  
and its tools to redirect the audio stream from the opened browser to the client using the service.