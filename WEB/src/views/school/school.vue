<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getSchools } from '../../services/school.service.js'
import { useNotification } from '../../composables/useNotification'
import TableCustom from '../../components/tables/TableCustom.vue'
import AppButton from '../../components/button/AppButton.vue'
import AppInput from '../../components/input/AppInput.vue'
import AppFilterBar from '../../components/common/AppFilterBar.vue'
import { debounce } from 'lodash-es'
import AppDialog from '@/components/dialogs/AppDialog.vue'
import AppForm from '@/components/forms/AppForm.vue'

const notify = useNotification()

const schools = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filters = reactive({ name: '' })

const columns = [
  { prop: 'name', label: 'ឈ្មោះសាលា' },
  { prop: 'abbreviation', label: 'អក្សរកាត់', width: 110 },
  { prop: 'director', label: 'នាយិកា' },
  { prop: 'phone', label: 'ទូរស័ព្ទ', width: 130 },
  { prop: 'email', label: 'អ៊ីមែល' },
  { prop: 'address', label: 'អាសយដ្ឋាន' },
]

// dialog / form state
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref(null)
const editingId = ref(null)

const form = reactive({
  name: '',
  abbreviation: '',
  director: '',
  phone: '',
  email: '',
  address: '',
  website: '',
  facebook: '',
})

const rules = {
  name: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះសាលា', trigger: 'blur' }],
  abbreviation: [{ required: true, message: 'សូមបញ្ចូលអក្សរកាត់', trigger: 'blur' }],
  director: [{ required: true, message: 'សូមបញ្ចូលឈ្មោះនាយក', trigger: 'blur' }],
  phone: [{ required: true, message: 'សូមបញ្ចូលលេខទូរស័ព្ទ', trigger: 'blur' }],
  email: [
    { required: true, message: 'សូមបញ្ចូលអ៊ីមែល', trigger: 'blur' },
    { type: 'email', message: 'អ៊ីមែលមិនត្រឹមត្រូវ', trigger: 'blur' },
  ],
  address: [{ required: true, message: 'សូមបញ្ចូលអាសយដ្ឋាន', trigger: 'blur' }],
}

async function fetchSchools() {
  loading.value = true
  try {
    const res = await getSchools({
      page: page.value,
      pageSize: pageSize.value,
      name: filters.name,
    })
    schools.value = res.data.data || []
    total.value = res.data.total || 0
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to load')
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(() => {
  page.value = 1
  fetchSchools()
}, 400)

function openCreate() {
  editingId.value = null
  dialogTitle.value = 'បង្កើតសាលាថ្មី'
  form.name = ''
  form.abbreviation = ''
  form.director = ''
  form.phone = ''
  form.email = ''
  form.address = ''
  form.website = ''
  form.facebook = ''
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  dialogTitle.value = 'កែប្រែសាលា'
  form.name = row.name
  form.abbreviation = row.abbreviation
  form.director = row.director
  form.phone = row.phone
  form.email = row.email
  form.address = row.address
  form.website = row.website
  form.facebook = row.facebook
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
}

async function handleSubmit() {
  submitting.value = true
  try {
    // TODO: call create/update service depending on editingId.value
    notify.success('រក្សាទុកបានជោគជ័យ')
    dialogVisible.value = false
    fetchSchools()
  } catch (e) {
    notify.error(e?.response?.data?.message || e.message || 'Failed to save')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchSchools()
})
</script>

<template>
  <div class="school-page">
    <AppFilterBar
      :fields="[{ slot: 'search', span: 18 }]"
      :action-span="6"
    >
      <template #search>
        <AppInput
          v-model="filters.name"
          placeholder="ស្វែងរកតាមឈ្មោះសាលា"
          clearable
          @input="debouncedFetch"
        />
      </template>
      <template #actions>
        <AppButton type="primary" icon="Plus" @click="openCreate">បង្កើតថ្មី</AppButton>
      </template>
    </AppFilterBar>

    <TableCustom
      show-index
      :data="schools"
      :columns="columns"
      :loading="loading"
      :total="total"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      @page-change="fetchSchools"
    >
      <template #actions="{ row }">
        <el-tooltip content="កែប្រែ" placement="top">
          <AppButton icon="Edit" circle size="small" type="warning" @click="openEdit(row)" />
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
        <AppInput
          v-model="form.name"
          placeholder="បញ្ចូលឈ្មោះសាលា"
          clearable
          prop="name"
          label="ឈ្មោះសាលា"
        />
        <AppInput
          v-model="form.abbreviation"
          placeholder="បញ្ចូលអក្សរកាត់"
          clearable
          prop="abbreviation"
          label="អក្សរកាត់"
        />
        <AppInput
          v-model="form.director"
          placeholder="បញ្ចូលឈ្មោះនាយក"
          clearable
          prop="director"
          label="នាយក"
        />
        <AppInput
          v-model="form.phone"
          placeholder="បញ្ចូលលេខទូរស័ព្ទ"
          clearable
          prop="phone"
          label="ទូរស័ព្ទ"
        />
        <AppInput
          v-model="form.email"
          placeholder="បញ្ចូលអ៊ីមែល"
          clearable
          prop="email"
          label="អ៊ីមែល"
        />
        <AppInput
          v-model="form.address"
          placeholder="បញ្ចូលអាសយដ្ឋាន"
          clearable
          prop="address"
          label="អាសយដ្ឋាន"
          type="textarea"
        />
        <AppInput
          v-model="form.website"
          placeholder="បញ្ចូលគេហទំព័រ (មិនចាំបាច់)"
          clearable
          prop="website"
          label="គេហទំព័រ"
        />
        <AppInput
          v-model="form.facebook"
          placeholder="បញ្ចូលតំណភ្ជាប់ Facebook (មិនចាំបាច់)"
          clearable
          prop="facebook"
          label="ហ្វេសប៊ុក"
        />
      </AppForm>
    </AppDialog>
  </div>
</template>

<style scoped>
.school-filters {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>