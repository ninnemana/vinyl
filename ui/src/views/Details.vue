<template>
    <div class='search'>
        <h1>Release Details</h1>
        <div>{{ this.$store.state.getResult }}</div>
    </div>
</template>

<script>
// @ is an alias to /src
// import Feed from '@/components/Feed.vue'

export default {
    name: 'Details',
    components: {}
}
</script>
