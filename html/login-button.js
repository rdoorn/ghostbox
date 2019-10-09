const template = document.createElement('template');
template.innerHTML = `
<style>
    :host {
    display: block;
    font-family: sans-serif;
    }

    .completed {
    text-decoration: line-through;
    }

    button {
    border: none;
    cursor: pointer;
    }
</style>
<button></button>
`;

class LoginButton extends HTMLElement {
  constructor() {
    super();
    this._shadowRoot = this.attachShadow({
      mode: 'open'
    });
    this._shadowRoot.appendChild(template.content.cloneNode(true));
    //this.$button = this._shadowRoot.querySelector('login-button');
    //this._shadowRoot.innerHTML = this.querySelector('label');
    //this._shadowRoot.innerHTML = this.attributes.label.value
    console.log(this.getRootNode())
    console.log(this.attributes.label.value)
    console.log(this._shadowRoot)

    this.$submitButton = this._shadowRoot.querySelector('button');
    this.$submitButton.innerHTML = this.attributes.label.value
    this.$submitButton.addEventListener('click', this._click.bind(this));

  }
  static get observedAttributes() {
    return ['label'];
  }
  attributeChangedCallback(name, oldVal, newVal) {
    this[name] = newVal;
    this.render();
  }
  /*this.addEventListener('click', () => {
    this.onClick('Hello from within the Custom Element');
  });*/
  render() {
    this.$button.innerHTML = this.label;
  }

  _click() {
    console.log('click!');

  }
  connectedCallback() {
    console.log('connected!');
  }

  disconnectedCallback() {
    console.log('disconnected!');
  }

  attributeChangedCallback(name, oldVal, newVal) {
    console.log(`Attribute: ${name} changed!`);
  }

  adoptedCallback() {
    console.log('adopted!');
  }

  static get observedAttributes() {
    return ['label'];
  }
}


customElements.define('login-button', LoginButton);