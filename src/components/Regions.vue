<template>
  <b-table :data="isEmpty ? [] : data.regions"
           striped
           :loading="isLoading"
           default-sort="price_monthly">
    <template slot-scope="props">
      <b-table-column field="slug" label="Name" sortable>
          {{ props.row.name }}
      </b-table-column>

      <b-table-column field="slug" label="Slug" sortable>
        <code>{{ props.row.slug }}</code>
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
      data: [],
      isLoading: true,
      isEmpty: false,
      errored: false
    }
  },
  created () {
    axios
      .get('/api/regions')
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
