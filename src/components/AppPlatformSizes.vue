<template>
  <section class="container">
    <div>
      <div class="level mb-4">
        <div class="level-left">
          <div class="level-item">
            <div style="width: 200px;">
              <label class="label">Filter by Tier</label>
              <b-select v-model="selectedTier" placeholder="All Tiers" expanded>
                <option value="">All Tiers</option>
                <option v-for="tier in tierTypes" :key="tier" :value="tier">
                  {{ formatTier(tier) }}
                </option>
              </b-select>
            </div>
          </div>
          <div class="level-item">
            <div style="width: 200px;">
              <label class="label">Filter by CPU Type</label>
              <b-select v-model="selectedCpuType" placeholder="All CPU Types" expanded>
                <option value="">All CPU Types</option>
                <option v-for="cpuType in cpuTypes" :key="cpuType" :value="cpuType">
                  {{ formatCpuType(cpuType) }}
                </option>
              </b-select>
            </div>
          </div>
          <div class="level-item">
            <div style="width: 200px; display: flex; align-items: flex-end; height: 73px;">
              <b-checkbox v-model="showDeprecated">
                Show deprecated
              </b-checkbox>
            </div>
          </div>
        </div>
      </div>

      <b-table :data="filteredSizes"
               striped
               :loading="isLoading"
               default-sort="usd_per_month"
               default-sort-direction="asc">
        <template slot-scope="props">
          <b-table-column field="slug" label="Slug" sortable>
            <code>{{ props.row.slug }}</code>
          </b-table-column>

          <b-table-column field="tier_slug" label="Tier" sortable>
              {{ formatTier(props.row.tier_slug) }}
          </b-table-column>

          <b-table-column field="cpu_type" label="CPU Type" sortable>
              {{ formatCpuType(props.row.cpu_type) }}
          </b-table-column>

          <b-table-column field="cpus" label="CPUs" sortable>
              {{ props.row.cpus }}
          </b-table-column>

          <b-table-column field="memory_bytes" label="Memory" sortable>
              {{ (props.row.memory_bytes / (1024 * 1024 * 1024)).toFixed(1) }} GB
          </b-table-column>

          <b-table-column field="bandwidth_allowance_gib" label="Bandwidth" sortable>
              {{ props.row.bandwidth_allowance_gib }} GiB
          </b-table-column>

          <b-table-column field="usd_per_month" label="Price Monthly" sortable :custom-sort="sortNumber">
              ${{ props.row.usd_per_month }}
          </b-table-column>

          <b-table-column field="usd_per_second" label="Price Per Second" sortable :custom-sort="sortNumber">
              ${{ props.row.usd_per_second }}
          </b-table-column>

          <b-table-column field="single_instance_only" label="Single Instance" sortable>
              {{ props.row.single_instance_only ? 'Yes' : 'No' }}
          </b-table-column>

          <b-table-column field="deprecation_intent" label="Deprecated" sortable>
              {{ props.row.deprecation_intent ? 'Yes' : 'No' }}
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
      selectedTier: '',
      selectedCpuType: '',
      showDeprecated: false
    }
  },
  computed: {
    filteredSizes () {
      if (this.isEmpty) return []
      if (!this.data.sizes) return []

      let filtered = this.data.sizes

      if (this.selectedTier) {
        filtered = filtered.filter(size => size.tier_slug === this.selectedTier)
      }

      if (this.selectedCpuType) {
        filtered = filtered.filter(size => size.cpu_type === this.selectedCpuType)
      }

      if (!this.showDeprecated) {
        filtered = filtered.filter(size => !size.deprecation_intent)
      }

      return filtered
    },
    tierTypes () {
      if (!this.data.sizes) return []

      // Extract unique tier types
      const uniqueTiers = new Set()
      this.data.sizes.forEach(size => {
        if (size.tier_slug) {
          uniqueTiers.add(size.tier_slug)
        }
      })

      return Array.from(uniqueTiers).sort()
    },
    cpuTypes () {
      if (!this.data.sizes) return []

      // Extract unique CPU types
      const uniqueCpuTypes = new Set()
      this.data.sizes.forEach(size => {
        if (size.cpu_type) {
          uniqueCpuTypes.add(size.cpu_type)
        }
      })

      return Array.from(uniqueCpuTypes).sort()
    }
  },
  methods: {
    sortNumber (a, b, key) {
      // Convert string values to numbers for proper sorting
      return parseFloat(a[key]) - parseFloat(b[key])
    },
    formatCpuType (cpuType) {
      if (!cpuType) return ''
      // Convert uppercase like "SHARED" to "Shared"
      return cpuType.charAt(0).toUpperCase() + cpuType.slice(1).toLowerCase()
    },
    formatTier (tier) {
      if (!tier) return ''
      // Convert lowercase like "basic" to "Basic"
      return tier.charAt(0).toUpperCase() + tier.slice(1)
    }
  },
  created () {
    axios
      .get('/api/apps/tiers/instance_sizes')
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
