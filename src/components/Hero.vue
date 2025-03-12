<template>
  <section class="hero is-medium is-primary">
    <div class="hero-head">
      <nav class="navbar">
        <div class="container">
          <div id="navbarMenuHeroA" class="navbar-menu">
            <div class="navbar-end">
              <span class="navbar-item">
                <a class="button is-primary is-inverted" href="https://docs.digitalocean.com/reference/api/api-reference/">
                  <span class="icon">
                    <i class="fab fa-digital-ocean"></i>
                  </span>
                  <span>API Documentation</span>
                </a>
              </span>
            </div>
          </div>
        </div>
      </nav>
    </div>

    <div class="hero-body">
      <div class="columns is-gapless">
        <div class="column is-one-fifth" />
        <div class="column is-one-fifth">
          <figure class="media-left">
            <p class="image is-128x128">
              <img src="../assets/logo-white.png" alt="Logo">
            </p>
          </figure>
        </div>
        <div class="column">
          <div class="container">
            <h1 class="title">
              DigitalOcean API Slugs
            </h1>
            <h2 class="subtitle">
              <code>{{ doctl }}</code>
            </h2>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'hero',
  data () {
    return {
      doctl: 'doctl compute size list'
    }
  },
  mounted () {
    this.$root.$on('tab', data => {
      if (data === 'sizes') {
        this.doctl = 'doctl compute size list'
      } else if (data === 'distro-images') {
        this.doctl = 'doctl compute image list-distribution'
      } else if (data === 'app-images') {
        this.doctl = 'doctl compute image list-application'
      } else if (data === 'k8s') {
        this.doctl = 'doctl k8s options versions'
      } else if (data === 'regions') {
        this.doctl = 'doctl compute region list'
      } else if (data === 'database-sizes') {
        this.doctl = 'doctl databases options slugs'
      } else if (data === 'database-versions') {
        this.doctl = 'doctl databases options versions'
      }
    })
  }
}
</script>
