<button
  on:click="{onclick}"
  class="c--button {outlined ? 'outlined' : ''} {text ? 'text' : ''} {red ? 'red' : ''} {green ? 'green' : ''}"
  title="{title}"
  disabled="{disabled}">
  <slot></slot>
</button>

<script>
  export let title
  export let onclick = () => {}

  export let red = false
  export let green = false

  export let disabled = false
  export let outlined = false
  export let text = false
  export let leading_icon = false
</script>

<style>
  .c--button {
    width: auto;
    min-width: 64px;
    height: 30px;
    background-color: #0277bd;
    -webkit-appearance: none;
    box-shadow: rgba(0, 0, 0, 0.2) 0 1px 2px -1px, rgba(0, 0, 0, 0.137255) 0 2px 3px 0, rgba(0, 0, 0, 0.117647) 0 1px 4px 0;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    color: rgb(255, 255, 255);
    outline-style: none;
    padding: 0 8px;
    font-size: 13px;
    margin: 2px;
    transition: background-color 200ms ease-out;
  }
  .c--button:hover {
    background-color: #0266ad;
    transition: background-color 160ms ease-in;
  }
  .c--button.red {
    background-color: #d32f2f;
  }
  .c--button.red:hover {
    background-color: #be2f2f;
  }
  .c--button.green {
    background-color: #66bb55;
  }
  .c--button.green:hover {
    background-color: #56a650;
  }

  .c--button.outlined {
    background-color: transparent;
    border: solid #0277bd 1px;
    color: #0277bd;
    box-shadow: none;
  }
  .c--button.outlined:hover {
    background-color: rgba(2, 119, 189, 0.05);
  }
  .c--button.outlined.red {
    border: solid #d32f2f 1px;
    color: #d32f2f;
  }
  .c--button.outlined.red:hover {
    background-color: rgba(190, 47, 47, 0.05);
  }
  .c--button.outlined.green {
    border: solid #66bb55 1px;
    color: #66bb55;
  }
  .c--button.outlined.green:hover {
    background-color: rgba(86, 166, 80, 0.05);
  }

  .c--button.text {
    background-color: transparent;
    color: #0277bd;
    box-shadow: none;
  }
  .c--button.text:hover {
    background-color: rgba(2, 119, 189, 0.05);
  }
  .c--button.text.red {
    color: #d32f2f;
  }
  .c--button.text.red:hover {
    background-color: rgba(190, 47, 47, 0.05);
  }
  .c--button.text.green {
    color: #66bb55;
  }
  .c--button.text.green:hover {
    background-color: rgba(86, 166, 80, 0.05);
  }

  .c--button:disabled {
    background: rgb(223, 223, 223);
    color: rgb(159, 159, 159);
    cursor: default;
    pointer-events: none;
    box-shadow: none;
  }
</style>