<template>
  <section class="container">
    <div>
      <div class="level mb-4">
        <div class="level-left">
          <div class="level-item">
            <div style="width: 200px;">
              <label class="label">Filter by Class</label>
              <b-select v-model="selectedClass" placeholder="All Classes" expanded>
                <option value="">All Classes</option>
                <option v-for="classType in classTypes" :key="classType" :value="classType">
                  {{ classType }}
                </option>
              </b-select>
            </div>
          </div>
        </div>
      </div>

      <b-table :data="filteredSizes"
               striped
               :loading="isLoading"
               default-sort="slug">
        <template slot-scope="props">
          <b-table-column field="description" label="Class" sortable>
              {{ props.row.description }}
          </b-table-column>

          <b-table-column field="slug" label="Slug" sortable>
            <code>{{ props.row.slug }}</code>
          </b-table-column>

          <b-table-column field="engines" label="Engines" sortable>
              {{ props.row.engines }}
          </b-table-column>

          <b-table-column field="nodes" label="Supported Node Count" sortable>
              {{ props.row.nodes.join(', ') }}
          </b-table-column>
        </template>
        <template slot="footer">
          <div class="has-text-right" v-if="data.retrieved_at">
              <span class="has-text-grey-light">Retrieved at: {{ data.retrieved_at }}</span>
          </div>
        </template>
        <template slot="empty">
          <section class="section is-medium">
            <div class="content has-text-grey has-text-centered">
              <div v-if="errored">
                <p>
                  <b-icon
                    pack="far"
                    icon="frown"
                    size="is-large">
                  </b-icon>
                </p>
                <p>Something went wrong here...</p>
              </div>
            </div>
          </section>
        </template>
      </b-table>
    </div>
  </section>
</template>

<script>
import axios from 'axios'
export default {
  data () {
    return {
      data: [],
      isLoading: true,
      isEmpty: false,
      errored: false,
      selectedClass: ''
    }
  },
  computed: {
    filteredSizes () {
      if (this.isEmpty) return []
      if (!this.data.sizes) return []

      if (!this.selectedClass) {
        return this.data.sizes
      }

      return this.data.sizes.filter(size => size.description === this.selectedClass)
    },
    classTypes () {
      if (!this.data.sizes) return []

      // Extract unique class types
      const uniqueClasses = new Set()
      this.data.sizes.forEach(size => {
        if (size.description) {
          uniqueClasses.add(size.description)
        }
      })

      return Array.from(uniqueClasses).sort()
    }
  },
  created () {
    axios
      .get('/api/databases/options')
      .then(response => {
        // Process the data from the API
        const options = response.data.options
        const sizes = []

        // Process each database type and its layouts
        Object.keys(options).forEach(engine => {
          const dbOptions = options[engine]

          // Process each layout (which contains supported node counts)
          dbOptions.layouts.forEach(layout => {
            const nodeCount = layout.num_nodes

            // Process each size in this layout
            layout.sizes.forEach(slug => {
              // Check if this size is already in our list
              let sizeIndex = sizes.findIndex(s => s.slug === slug)

              if (sizeIndex === -1) {
                sizes.push({
                  slug: slug,
                  engines: [engine],
                  nodes: [nodeCount],
                  description: this.getClassFromSlug(slug)
                })
              } else {
                // This size already exists in our list, so we need to update it
                const size = sizes[sizeIndex]

                // Add the engine if it's not already there
                if (!size.engines.includes(engine)) {
                  size.engines.push(engine)
                }

                // Add the node count if it's not already there
                if (!size.nodes.includes(nodeCount)) {
                  size.nodes.push(nodeCount)
                  size.nodes.sort((a, b) => a - b) // Sort numerically
                }
              }
            })
          })
        })

        // Sort engines for each size
        sizes.forEach(size => {
          size.engines.sort()
          size.engines = size.engines.join(', ')
        })

        // Now use the API to get the class descriptions if needed
        return axios.get('/api/sizes').then(sizesResponse => {
          const dropletSizes = sizesResponse.data.sizes

          // Update our database sizes with class descriptions from the Droplet sizes if available
          sizes.forEach(dbSize => {
            const matchingSize = dropletSizes.find(s => s.slug === dbSize.slug)
            if (matchingSize && matchingSize.description) {
              dbSize.description = matchingSize.description
            }
          })

          return {
            sizes: sizes,
            retrieved_at: response.data.retrieved_at
          }
        })
      })
      .then(processedData => {
        this.data = processedData
      })
      .catch(error => {
        console.log(error)
        this.isEmpty = true
        this.errored = true
      })
      .finally(() => { this.isLoading = false })
  },
  methods: {
    getClassFromSlug (slug) {
      // Extract class from slug if possible
      if (slug.startsWith('db-')) {
        return 'Database'
      } else if (slug.startsWith('gd-')) {
        return 'General Purpose'
      } else if (slug.startsWith('m-') || slug.startsWith('m3-')) {
        return 'Memory Optimized'
      } else if (slug.startsWith('so-') || slug.startsWith('so1_5-')) {
        return 'Storage Optimized'
      }
      return 'Unknown'
    }
  }
}
</script>
