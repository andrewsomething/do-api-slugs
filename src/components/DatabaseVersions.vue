<template>
  <b-table :data="isEmpty ? [] : formattedData"
           striped
           :loading="isLoading"
           default-sort="engine">
    <template slot-scope="props">
      <b-table-column field="engine" label="Engine" sortable>
        <code>{{ props.row.engine }}</code>
      </b-table-column>

      <b-table-column field="versions" label="Versions" sortable>
        {{ props.row.versions }}
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
</template>

<script>
import axios from 'axios'

export default {
  data () {
    return {
      data: {
        options: {}
      },
      isLoading: true,
      isEmpty: false,
      errored: false
    }
  },
  computed: {
    formattedData () {
      if (!this.data.options) {
        return []
      }

      return Object.entries(this.data.options).map(([engine, data]) => {
        return {
          engine: engine,
          versions: data.versions.join(', ')
        }
      })
    }
  },
  created () {
    axios
      .get('/api/databases/options')
      .then(response => {
        this.data = response.data
      })
      .catch(error => {
        console.log(error)
        this.isEmpty = true
        this.errored = true
      })
      .finally(() => { this.isLoading = false })
  }
}
</script>
