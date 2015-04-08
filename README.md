# Triangle

Send commands to your currently playing audio source, local music player like MPD or even your SoundCloud in your browser.  
Triangle is programmed in Go, while the browser adapters are in JavaScript.

## Status
Currently Triangle is in early development stage, PoC you could say.  
What we have:

* RPC Server that can handle the toggle command on the most recent playing adapter.
* MPD adapter.
* Chromium (SoundCloud) adapter.

## TODO
* Improve code quality.
* Write tests.
* Add overlay notifications on events(source/track changed).
* Implement the prev and next commands.
* Implement YouTube support on the Chromium adapter.
* Add Firefox adatper.
* Add VLC/WinAmp/AIMP/WMP(etc...) adapters.
