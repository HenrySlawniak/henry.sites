'use strict';

document.addEventListener("DOMContentLoaded", function () {
  var vid = document.getElementById("video");
  vid.volume = 0.2;
  vid.play();
  vid.onclick = function () {
    if (vid.paused) {
      vid.play();
    } else {
      vid.pause();
    }
  }
});
