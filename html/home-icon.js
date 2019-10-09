const template = document.createElement('template')
template.innerHTML = `
  <style>
  img {
    width:50px;
    height:50px;
    border:0px;
  }
  a {
    font-size: 30px;
    text-decoration: none;
    color:#aaa;
    line-height: 50px;
  }
  </style>
  <img src=""/><a href=""></a>
`;

class HomeIcon extends HTMLElement {
  constructor() {
    super();
    this._shadowRoot = this.attachShadow({
      mode: 'open'
    });

    this.authenticated = (window.sessionStorage.accessToken != null)


    this._shadowRoot.appendChild(template.content.cloneNode(true));


    console.log(this.getRootNode())
    console.log(this.attributes)
    console.log(this._shadowRoot)

    this.$a = this._shadowRoot.querySelector('a');
    this.$img = this._shadowRoot.querySelector('img');
    //this.$submitButton.innerHTML = this.attributes.label.value
    //this.$submitButton.addEventListener('click', this._click.bind(this));
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
    //this.$button.innerHTML = this.label;
    this.$img.src = this.getAttribute("image")
    this.$a.innerHTML = this.getAttribute("name")

    if (this.authenticated) {
      this.$a.href = this.getAttribute("userhome")
    } else {
      this.$a.href = this.getAttribute("home")
    }
  }

  _click() {
    console.log('click!');

  }
  connectedCallback() {
    console.log('connected!');
    this.render()
  }

  disconnectedCallback() {
    console.log('disconnected!');
  }

  attributeChangedCallback(name, oldVal, newVal) {
    console.log(`
Attribute: $ {
  name
}
changed!`);
  }

  adoptedCallback() {
    console.log('adopted!');
  }

  static get observedAttributes() {
    return ['label'];
  }

}

customElements.define('home-icon', HomeIcon);