<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getFeediscountGroup,

  createFeediscountGroup,
  updateFeediscountGroup,
  toggleFeediscountGroup
 } from '../../services/feediscountgroup.service.js'
import { useNotification } from '../../composables/useNotification'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppButton from '../../components/button/AppButton.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppSelect from '../../components/common/AppSelect.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'
import { debounce } from 'lodash-es'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppForm from '@/components/forms/AppForm.vue'

const notify = useNotification()

const feeDiscountGroups = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const statusOptions = [
  { label: 'សកម្ម', value: 1 },
  { label: 'អសកម្ម', value: 0 },
]

const discountTypeOptions = [
  { label: 'ភាគរយ', value: 'percentage' },
  { label: 'ចំនួនទឹកប្រាក់', value: 'amount' },
]

const filters = reactive({ name: '', active: null })

const columns = [
  { prop: 'code', label: 'កូដ', slot: 'code' },
  { prop: 'name', label: 'ឈ្មោះ' },
  { prop: 'discount_type', label: 'ប្រភេទបញ្ចុះតម្លៃ', slot: 'discountType' },
  { prop: 'discount_value', label: 'តម្លៃបញ្ចុះ', slot: 'discountValue' },
  { prop: 'description', label: 'ការពិពណ៌នា' },
  { prop: 'active', label: 'ស្ថានភាព', slot: 'isActive', width: 100 },
]

// dialog / form state
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingId = ref(null)

const form = reactive({
  code: '',
  name: '',
  discount_type: 'percentage',
  discount_percentage: 0,
  discount_amount: 0,
  description: '',
})

const rules = {
  code: [{ required: true, message: 'សូមបញ្ចូលកូដ', trigger: 'blur' }],
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះ', trigger: 'blur' }],
  discount_type: [{ required: true, message: 'សូមជ្រើសរើសប្រភេទបញ្ចុះតម្លៃ', trigger: 'change' }],
}

async function fetchFeeDiscountGroups() {
  loading.value = true
  try {
    const res = await getFeediscountGroup({
      page: page.value,
      pageSize: pageSize.value,
      name: filters.name,
      active: filters.active,
    })
    feeDiscountGroups.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1
  fetchFeeDiscountGroups()
}, 400)

function openCreate() {
  editingId.value = null
  dialogTitle.value = 'បង្កើតក្រុមបញ្ចុះតម្លៃថ្មី'
  form.code = ''
  form.name = ''
  form.discount_type = 'percentage'
  form.discount_percentage = 0
  form.discount_amount = 0
  form.description = ''
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.uuid ?? row.id
  dialogTitle.value = 'កែប្រែក្រុមបញ្ចុះតម្លៃ'
  form.code = row.code
  form.name = row.name
  form.discount_type = row.discount_type
  form.discount_percentage = row.discount_percentage
  form.discount_amount = row.discount_amount
  form.description = row.description
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (form.discount_type === 'percentage') {
      form.discount_amount = 0
    } else {
      form.discount_percentage = 0
    }

    if (editingId.value) {
      await updateFeediscountGroup(editingId.value, { ...form })
    } else {
      await createFeediscountGroup({ ...form })
    }

    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchFeeDiscountGroups()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(row) {
  try {
    await toggleFeediscountGroup(row.uuid ?? row.id)
    notify.success('ធ្វើបច្ចុប្បន្នភាពស្ថានភាពបានជោគជ័យ')
    fetchFeeDiscountGroups()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to toggle status')
  }
}

onMounted(() => {
  fetchFeeDiscountGroups()
})
</script>

<template>
  <div class="category-page">
    <AppFilterBar
      :fields="[
        { slot: 'search', span: 10 },
        { slot: 'create', span: 5 },
      ]"
    >
      <template #search>
        <AppInput
          v-model="filters.name"
          placeholder="ស្វែងរកតាមឈ្មោះ"
          clearable
          @input="debouncedFetch"
        />
      </template>
      <!-- <AppSelect
        v-model="filters.active"
        :options="statusOptions"
        placeholder="ស្ថានភាព"
        clearable
        @change="debouncedFetch"
      /> -->
      <template #create>
        <AppButton type="primary" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      show-index
      :data="feeDiscountGroups"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchFeeDiscountGroups"
    >
      <template #discountType="{ row }">
        {{ row.discount_type === 'percentage' ? 'ភាគរយ' : 'ចំនួនទឹកប្រាក់' }}
      </template>

      <template #code="{ row }">
        <el-text tag="b" style="color: darkcyan">
          {{ row.code }}
        </el-text>
      </template>

      <template #discountValue="{ row }">
        <el-text tag="b" style="color: crimson">
          {{
            row.discount_type === 'percentage'
              ? `${row.discount_percentage}%`
              : `${row.discount_amount}$`
          }}
        </el-text>
      </template>

      <template #isActive="{ row }">
        <el-tag :type="row.active ? 'success' : 'danger'">
          {{ row.active ? 'សកម្ម' : 'អសកម្ម' }}
        </el-tag>
      </template>

      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton icon="Edit" circle size="small" type="warning" @click="openEdit(row)" />
        </el-tooltip>
        <el-tooltip content="បិទ/បេីក" placement="top">
          <AppButton
            :icon="row.active ? 'CircleClose' : 'CircleCheck'"
            circle
            size="small"
            type="danger"
            @click="toggleStatus(row)"
          />
        </el-tooltip>
      </template>
    </TableCustom>

    <AppDialog v-model:visible="dialogVisible" :title="dialogTitle" :showDefaultFooter="false">
      <AppForm
        ref="formRef"
        :model="form"
        :rules="rules"
        :loading="submitting"
        :show-actions="true"
        @submit="handleSubmit"
        @reset="closeDialog"
        submitText="រក្សាទុក"
        resetText="ចាកចេញ"
      >
        <AppInput v-model="form.code" placeholder="បញ្ចូលកូដ" clearable prop="code" label="កូដ" />
        <AppInput
          v-model="form.name"
          placeholder="បញ្ចូលឈ្មោះ"
          clearable
          prop="name"
          label="ឈ្មោះ"
        />

        <AppSelect
          v-model="form.discount_type"
          :options="discountTypeOptions"
          prop="discount_type"
          label="ប្រភេទបញ្ចុះតម្លៃ"
        />

        <AppInput
          v-if="form.discount_type === 'percentage'"
          v-model.number="form.discount_percentage"
          type="number"
          placeholder="បញ្ចូលភាគរយ"
          prop="discount_percentage"
          label="ភាគរយបញ្ចុះតម្លៃ (%)"
        />

        <AppInput
          v-else
          v-model.number="form.discount_amount"
          type="number"
          placeholder="បញ្ចូលចំនួនទឹកប្រាក់"
          prop="discount_amount"
          label="ចំនួនទឹកប្រាក់បញ្ចុះតម្លៃ"
        />

        <AppInput
          v-model="form.description"
          placeholder="បញ្ចូលការពិពណ៌នា"
          clearable
          prop="description"
          label="ការពិពណ៌នា"
          type="textarea"
        />
      </AppForm>
    </AppDialog>
  </div>
</template>

<style scoped>
.category-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>
