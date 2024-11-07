import App from './App.svelte'
import './global.css'
import { mount } from "svelte";

export default mount(App, { target: document.body })
