import React, { useState, useCallback, useContext, useMemo } from 'react'
import { StyleSheet, View, Text, TouchableOpacity } from 'react-native'

import FormLayout from '../../Components/FormLayout'
import DataField from '../../Components/DataField'

import { CreateNewRole } from '../../Services/RoleService'

import { StoreContext } from '../../App'

const CreateRole = ({ navigation }) => {
  const { setLoading, currentProject } = useContext(StoreContext)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')

  const [create, setCreate] = useState(false)
  const [add, setAdd] = useState(false)
  const [assign, setAssign] = useState(false)

  const permissions = useMemo(() => {
    let temp = ''
    if (create) {
      temp += '2'
    }
    if (assign) {
      temp += '3'
    }
    if (add) {
      temp += '4'
    }
    return temp.split('').join(',')
  }, [add, assign, create])

  const onSubmit = useCallback(async () => {
    await setLoading(true)
    const res = await CreateNewRole({
      name,
      description,
      permissions,
      project_id: currentProject.ID,
    })
    await setLoading(false)
    if (res) {
      navigation.goBack()
    }
  }, [
    currentProject.ID,
    description,
    name,
    navigation,
    permissions,
    setLoading,
  ])

  return (
    <View style={styles.container}>
      <FormLayout
        title="CREATE NEW ROLE"
        onSubmit={onSubmit}
        disabled={
          permissions.length < 1 || name.length < 1 || description.length < 1
        }
      >
        <DataField label="Name" value={name} onChange={setName} />
        <DataField
          label="Description"
          value={description}
          onChange={setDescription}
        />
      </FormLayout>
      <View style={styles.permission}>
        <Text style={styles.bold}>Select permission:</Text>
        <TouchableOpacity
          onPress={() => setCreate(!create)}
          style={styles.selection}
        >
          <View style={styles.radio}>
            {create && <View style={styles.selected} />}
          </View>
          <Text>Create/Delete task</Text>
        </TouchableOpacity>

        <TouchableOpacity onPress={() => setAdd(!add)} style={styles.selection}>
          <View style={styles.radio}>
            {add && <View style={styles.selected} />}
          </View>
          <Text>Assign/Unassign member from task</Text>
        </TouchableOpacity>

        <TouchableOpacity
          onPress={() => setAssign(!assign)}
          style={styles.selection}
        >
          <View style={styles.radio}>
            {assign && <View style={styles.selected} />}
          </View>
          <Text>Add/Remove member from project</Text>
        </TouchableOpacity>
      </View>
    </View>
  )
}

export default CreateRole

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  permission: {
    flex: 1,
    paddingHorizontal: 32,
  },
  bold: {
    fontWeight: 'bold',
  },
  selection: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: 16,
  },
  radio: {
    width: 20,
    height: 20,
    borderRadius: 20,
    padding: 2,
    borderColor: 'firebrick',
    borderWidth: 2,
    marginRight: 8,
  },
  selected: {
    backgroundColor: 'firebrick',
    flex: 1,
    borderRadius: 20,
  },
})
