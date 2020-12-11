import React, { useState, useCallback, useContext } from 'react'
import { StyleSheet, View } from 'react-native'
import { StoreContext } from '../../App'

import FormLayout from '../../Components/FormLayout'
import DataField from '../../Components/DataField'

import { UpdateProjectInfo } from '../../Services/ProjectService'

const ProjectInfo = () => {
  const { setLoading, currentProject, setCurrentProject } = useContext(
    StoreContext,
  )
  const [name, setName] = useState(currentProject.name || '')
  const [description, setDescription] = useState(
    currentProject.description || '',
  )

  const onSubmit = useCallback(async () => {
    await setLoading(true)
    const res = await UpdateProjectInfo({
      name,
      description,
      ID: currentProject.ID,
    })
    await setCurrentProject(res)
    await setLoading(false)
  }, [setLoading, name, description, currentProject, setCurrentProject])
  return (
    <View style={styles.flex1}>
      <View style={styles.container}>
        <FormLayout title="UPDATE PROJECT INFO" onSubmit={onSubmit}>
          <DataField label="Name" value={name} onChange={setName} />
          <DataField
            label="Description"
            value={description}
            onChange={setDescription}
          />
        </FormLayout>
      </View>
    </View>
  )
}

export default ProjectInfo

const styles = StyleSheet.create({
  container: {
    flex: 0.5,
  },
  flex1: {
    flex: 1,
    backgroundColor: 'black',
  },
})
