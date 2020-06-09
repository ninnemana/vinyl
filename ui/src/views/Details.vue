<template>
    <div class='search'>
        <Release />
    </div>
</template>

<script>
// @ is an alias to /src
import Release from '@/components/Release.vue'

export default {
    name: 'Details',
    components: {
        Release,
    },
};
</script>
