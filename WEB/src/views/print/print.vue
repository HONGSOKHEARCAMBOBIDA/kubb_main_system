<template>
  <PrintableTable
    title="Student List"
    :columns="columns"
    :rows="students"
  />
</template>

<script setup>
import { ref } from "vue"
import PrintableTable from "../../layouts/components/PrintableTable.vue";

const columns = [
  { key: "id", label: "Student ID" },
  { key: "name", label: "Name" },
  { key: "gender", label: "Gender" },
  { key: "grade", label: "Grade" },
  { key: "score", label: "Score" }
]

const students = ref([
  { id: 1001, name: "Dara", gender: "Male", grade: "10A", score: 85 },
  { id: 1002, name: "Sokha", gender: "Female", grade: "10A", score: 92 },
  { id: 1003, name: "Vanna", gender: "Male", grade: "10B", score: 78 }
])
</script>