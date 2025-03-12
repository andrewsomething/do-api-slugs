<template>
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
             default-sort="price_monthly">
      <template slot-scope="props">
        <b-table-column field="description" label="Class" sortable>
            {{ props.row.description }}
        </b-table-column>

        <b-table-column field="slug" label="Slug" sortable>
          <code>{{ props.row.slug }}</code>
        </b-table-column>

        <b-table-column field="memory" label="RAM" sortable>
            {{ props.row.memory | mbToGb }} GB
        </b-table-column>

        <b-table-column field="vcpus" label="CPU" sortable>
            {{ props.row.vcpus }}
        </b-table-column>

        <b-table-column field="disk" label="Disk" sortable>
            {{ props.row.disk }} GB
        </b-table-column>

        <b-table-column field="transfer" label="Transfer" sortable>
            {{ props.row.transfer }} TB
        </b-table-column>

        <b-table-column field="price_monthly" label="Price Monthly" sortable>
            ${{ props.row.price_monthly }}
        </b-table-column>

        <b-table-column field="price_hourly" label="Price Hourly" sortable>
            ${{ props.row.price_hourly }}
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
  filters: {
    mbToGb: function (value) {
      return value / 1024
    }
  },
  created () {
    axios
      .get('/api/sizes')
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
