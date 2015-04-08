function main() {
  console.log("Triangle extension is loaded");
  connectToTriangle();
}

function connectToTriangle() {
  console.log("Trying to connect to Triangle");
  var client = new WebSocket("ws://127.0.0.1:8080");
  client.onopen = function() {
    console.log("Connected to Triangle");

    // Send currently playing state
    isSoundCloudPlaying(function(result) {
      client.send(result);
    });

    registerOnPlayerChange(client);
  };

  client.onclose = function() {
    console.log("Disconnected from Triangle");
    setTimeout(connectToTriangle, 2000);
  };

  client.onmessage = function(evt) {
    console.log("Received message from Triangle", evt.data);
    if (evt.data == "toggle") {
      togglePlay();
    } else if (evt.data == "info") {
      isSoundCloudPlaying(function(result) {
        client.send(result);
      });
    }
  };
}

function injectedToSoundCloud() {
  var playButton = document.querySelector('button.playControl');
  if (window.triangleLoaded) {
    console.log("Triangle is already injected");
  } else {
    console.log("Injecting Triangle");
    function onClick(e) {
      var className = e.target.className;
      if (className.indexOf("playControl") > -1 || 
          className.indexOf("sc-button-play") > -1) {
        var isPlaying = playButton.className.indexOf("playing") > -1;
        chrome.runtime.sendMessage(isPlaying);
      }
    }
    document.body.addEventListener("click", onClick);
    window.triangleLoaded = true;
  }
}

function registerOnPlayerChange(client) {
  chrome.runtime.onMessage.addListener(function(message) {
    console.log("Sending to Triangle", message);
    client.send(message);
  });
  function onSoundCloudLoad(tabId) {
    chrome.tabs.get(tabId, function(tab) {
      if (tab.url.indexOf("https://soundcloud.com") > -1) {
        chrome.tabs.executeScript(
            tabId,
            {code: "(" + injectedToSoundCloud.toString() + ")()"}
            );
      }
    });
  };
  chrome.tabs.query({url: "https://soundcloud.com/*"}, function(tabs) {
    tabs.forEach(function(tab) {
      onSoundCloudLoad(tab.id);
    });
  });
  chrome.tabs.onUpdated.addListener(onSoundCloudLoad);
}

function isSoundCloudPlaying(callback) {
  chrome.tabs.query({url: "https://soundcloud.com/*"}, function(tabs) {
    if (tabs.length < 1) {
      // No SoundCloud tab found
      callback(false);
      return;
    }
    chrome.tabs.executeScript(
        tabs[0].id,
        {code: "!!document.querySelector('button.playControl.playing')"},
        function(results) {
          callback(results[0]);
        }
        );
  });
}

function togglePlay() {
  chrome.tabs.query({url: "https://soundcloud.com/*"}, function(tabs) {
    if (tabs.length < 1) {
      console.log("No SoundCloud tab");
      return;
    }
    chrome.tabs.executeScript(
        tabs[0].id,
        {code: "document.querySelector('button.playControl').click()"}
        );
  });
}

main();
