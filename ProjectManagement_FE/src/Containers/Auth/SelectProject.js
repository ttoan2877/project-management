import React, { useEffect, useCallback, useContext, useState } from 'react'
import {
  StyleSheet,
  Text,
  SafeAreaView,
  TouchableOpacity,
  View,
  FlatList,
} from 'react-native'

import AsyncStorage from '@react-native-async-storage/async-storage'

import { FetchProject, GetProjectById } from '../../Services/ProjectService'

import { StoreContext } from '../../App'

const ProjectItem = ({ data, onPress }) => {
  return (
    <TouchableOpacity onPress={() => onPress(data.ID)} style={styles.item}>
      <Text>
        ID: <Text style={styles.bold}>{data.ID}</Text>
      </Text>
      <Text>
        Name: <Text style={styles.bold}>{data.name}</Text>
      </Text>
      <Text>
        Description: <Text style={styles.bold}>{data.description}</Text>
      </Text>
    </TouchableOpacity>
  )
}

const SelectProject = ({ navigation }) => {
  const { setLoading, setCurrentProject } = useContext(StoreContext)
  const [projects, setProjects] = useState([])

  useEffect(() => {
    const getProject = async () => {
      await setLoading(true)
      const data = await FetchProject()
      await setLoading(false)
      data && setProjects(data)
    }
    getProject()
  }, [setLoading])

  const onPressItem = useCallback(
    async (ID) => {
      await setLoading(true)
      const project = await GetProjectById(ID)
      await setCurrentProject(project)
      await setLoading(false)
      return navigation.navigate('Main')
    },
    [navigation, setCurrentProject, setLoading],
  )

  const onLogout = useCallback(async () => {
    await AsyncStorage.removeItem('@accessToken')
    navigation.navigate('Login')
  }, [navigation])

  return (
    <View style={styles.flex1}>
      <SafeAreaView style={styles.container}>
        <Text style={styles.title}>JOIN YOUR PROJECT</Text>
        <FlatList
          style={styles.list}
          data={projects}
          renderItem={({ item, index }) => (
            <ProjectItem
              key={index.toString()}
              data={item}
              onPress={onPressItem}
            />
          )}
          ListEmptyComponent={
            <Text style={styles.empty}>Please create a project!</Text>
          }
        />
        <View style={styles.btnWrapper}>
          <TouchableOpacity
            onPress={onLogout}
            style={[styles.btn, styles.btnLogout]}
          >
            <Text style={[styles.label, styles.labelLogout]}>Log out</Text>
          </TouchableOpacity>
          <TouchableOpacity
            onPress={() => navigation.navigate('CreateProject')}
            style={styles.btn}
          >
            <Text style={styles.label}>New project</Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    </View>
  )
}

export default SelectProject

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginVertical: 64,
    marginHorizontal: 32,
    borderColor: 'firebrick',
    borderWidth: 4,
    borderRadius: 16,
    backgroundColor: 'black',
  },
  list: {
    flex: 1,
    marginVertical: 16,
  },
  title: {
    color: 'firebrick',
    fontWeight: 'bold',
    marginTop: 16,
    alignSelf: 'center',
  },
  btn: {
    flex: 1,
    backgroundColor: 'firebrick',
    paddingVertical: 12,
    borderRadius: 32,
    marginHorizontal: 8,
    justifyContent: 'center',
    alignItems: 'center',
  },
  label: {
    fontWeight: 'bold',
    color: 'azure',
  },
  btnWrapper: {
    marginBottom: 16,
    flexDirection: 'row',
  },
  btnLogout: {
    backgroundColor: 'azure',
    borderColor: 'firebrick',
    borderWidth: 4,
  },
  labelLogout: {
    color: 'black',
  },
  item: {
    padding: 16,
    borderWidth: 2,
    margin: 8,
    borderColor: 'firebrick',
    borderRadius: 8,
    backgroundColor: 'azure',
  },
  empty: {
    alignSelf: 'center',
  },
  bold: {
    fontWeight: 'bold',
  },
  flex1: {
    flex: 1,
    backgroundColor: 'black',
  },
})
