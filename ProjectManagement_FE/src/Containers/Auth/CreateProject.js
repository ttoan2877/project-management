import React, { useState, useCallback, useContext } from 'react'

import {
  StyleSheet,
  Text,
  View,
  TextInput,
  TouchableOpacity,
} from 'react-native'

import { CreateNewProject, GetProjectById } from '../../Services/ProjectService'

import { StoreContext } from '../../App'

const CreateProject = ({ navigation }) => {
  const { setLoading, setCurrentProject } = useContext(StoreContext)

  const [name, setName] = useState('')

  const [description, setDescription] = useState('')

  const onSubmit = useCallback(async () => {
    await setLoading(true)
    const createdProject = await CreateNewProject({ name, description })
    const project = await GetProjectById(createdProject.ID)
    await setCurrentProject(project)
    await setLoading(false)
    return navigation.navigate('Main')
  }, [description, name, navigation, setCurrentProject, setLoading])

  return (
    <View style={styles.container}>
      <Text style={styles.title}>CREATE NEW PROJECT</Text>
      <Text style={styles.label}>PROJECT NAME:</Text>
      <TextInput
        style={[styles.input, styles.h48]}
        value={name}
        onChangeText={(val) => setName(val)}
      />
      <Text style={styles.label}>PROJECT DESCRIPTION:</Text>
      <TextInput
        multiline
        style={[styles.input, styles.h64]}
        value={description}
        onChangeText={(val) => setDescription(val)}
      />
      <TouchableOpacity onPress={onSubmit} style={styles.btn}>
        <Text style={styles.btnLabel}>CREATE</Text>
      </TouchableOpacity>
      <TouchableOpacity onPress={() => navigation.goBack()}>
        <Text style={styles.goBack}>Back to select</Text>
      </TouchableOpacity>
    </View>
  )
}

export default CreateProject

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
    backgroundColor: 'black',
  },
  input: {
    width: '100%',
    borderWidth: 4,
    borderColor: 'firebrick',
    margin: 16,
    paddingHorizontal: 32,
    fontSize: 16,
    textAlign: 'center',
    backgroundColor: 'azure',
  },
  h48: {
    height: 48,
    borderRadius: 16,
    fontWeight: 'bold',
  },
  h64: {
    minHeight: 64,
    borderRadius: 16,
    paddingVertical: 16,
  },
  label: {
    fontSize: 12,
    fontWeight: 'bold',
    color: 'azure',
  },
  title: {
    marginBottom: 64,
    fontSize: 16,
    fontWeight: 'bold',
    color: 'firebrick',
  },
  btn: {
    width: '80%',
    backgroundColor: 'firebrick',
    padding: 16,
    alignItems: 'center',
    borderRadius: 64,
    marginTop: 32,
  },
  btnLabel: {
    color: 'azure',
    fontWeight: 'bold',
  },
  goBack: {
    color: 'firebrick',
    fontWeight: 'bold',
    marginTop: 16,
    fontSize: 16,
  },
})
