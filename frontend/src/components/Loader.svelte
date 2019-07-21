<div class="loader--wrapper"></div>

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
  .loader--wrapper:after {
    content: " ";
    display: block;
    width: 46px;
    height: 46px;
    margin: 1px;
    border-radius: 50%;
    border: 5px solid #4a57ff;
    border-color: #4a57ff transparent #4a57ff transparent;
    animation: loader--wrapper-animation 1.2s linear infinite;
  }

  @keyframes loader--wrapper-animation {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
</style>