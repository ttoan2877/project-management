import React, { useState, useCallback, useContext } from 'react'
import { StyleSheet, View, ScrollView } from 'react-native'
import DataField from '../../Components/DataField'
import FormLayout from '../../Components/FormLayout'

import { CreateTask } from '../../Services/TaskService'

import { StoreContext } from '../../App'

const FormTask = ({ navigation }) => {
  const { setLoading, currentProject } = useContext(StoreContext)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')

  const label = 'CREATE NEW TASK'

  const onSubmit = useCallback(async () => {
    await setLoading(true)
    const res = await CreateTask({
      name,
      description,
      project_id: currentProject.ID,
    })
    if (res) {
    }
    await setLoading(false)
  }, [currentProject, description, name, setLoading])

  return (
    <ScrollView style={styles.container}>
      <FormLayout title={label} onSubmit={onSubmit}>
        <DataField label="Name" value={name} onChange={setName} />
        <DataField
          label="Description"
          value={description}
          onChange={setDescription}
        />
      </FormLayout>
      <View style={styles.userWrapper}></View>
      <View style={styles.subWrapper}></View>
    </ScrollView>
  )
}

export default FormTask

const styles = StyleSheet.create({
  container: {
    backgroundColor: 'black',
  },
  userWrapper: {
    flex: 1,
  },
  subWrapper: {
    flex: 1,
  },
})
