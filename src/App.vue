<template>
  <div id="app">
    <Hero />
    <TabbedView />
    <Footer />
  </div>
</template>

<script>
import Hero from './components/Hero.vue'
import TabbedView from './components/TabbedView.vue'
import Footer from './components/Footer.vue'
export default {
  name: 'app',
  components: {
    Hero,
    TabbedView,
    Footer
  }
}
</script>

<style lang="scss">
  @import "~bulma/sass/utilities/_all";

  $primary: #0080FF;
  $primary-invert: findColorInvert($primary);
  $twitter: #4099ff;
  $twitter-invert: findColorInvert($twitter);
  $colors: (
      "white": ($white, $black),
      "black": ($black, $white),
      "light": ($light, $light-invert),
      "dark": ($dark, $dark-invert),
      "primary": ($primary, $primary-invert),
      "info": ($info, $info-invert),
      "success": ($success, $success-invert),
      "warning": ($warning, $warning-invert),
      "danger": ($danger, $danger-invert),
      "twitter": ($twitter, $twitter-invert)
  );
  // Links
  $link: $primary;
  $link-invert: $primary-invert;
  $link-focus-border: $primary;

  $footer-background-color: #010E28;
  .fa-lg {
    padding: 5px;
  }
  .footer.a {
    color: $red;
  }

  @import "~bulma";
  @import "~buefy/src/scss/buefy";
</style>
