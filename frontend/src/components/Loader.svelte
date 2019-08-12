<div class="loader--wrapper">
  <div class="o first"></div>
  <div class="o second"></div>
  <div class="o third"></div>
</div>

<style>
  .loader--wrapper {
    display: inline-block;
    width: 64px;
    height: 64px;
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    margin: auto;
  }

  .o {
    box-sizing: border-box;
    position: absolute;
    border-radius: 50%;
    width: 60px;
    height: 60px;
    border: 2px solid rgba(0, 0, 60, .1);
  }

  .first {
    animation: loader--animation 1.1s infinite cubic-bezier(0, 0, 1, 1);
    border-left: 3px solid rgba(0, 0, 60, .7);
  }

  .second {
    animation: loader--animation 0.6s infinite cubic-bezier(0, 0, 1, 1);
    border-left: 3px solid rgba(255, 0, 0, .7);
  }

  .third {
    animation: loader--animation 0.9s infinite cubic-bezier(0, 0, 1, 1);
    border-left: 3px solid rgba(100, 100, 255, .7);
  }

  @keyframes loader--animation {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
</style>