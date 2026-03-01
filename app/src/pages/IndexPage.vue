<template>
  <q-page padding>
    <ApexChart
      type="donut"
      height="350"
    />
  </q-page>
</template>

<script setup lang="ts">
    import ApexChart from 'vue3-apexcharts';

    const chartOptions = {
        series: [44, 55, 13],
        labels: ['A', 'B', 'C']
    };

    const chartSeries = [
        {
            name: 'sales',
            data: [30,40,35,50,49,60,70,91,125]
        }
    ];
</script>