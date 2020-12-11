import React, { useEffect, useCallback, useContext, useState } from 'react'

import { View, Text, TouchableOpacity, StyleSheet } from 'react-native'

import { FetchTaskInProject } from '../../Services/TaskService'

import { StoreContext } from '../../App'

import CommonLayout from '../../Components/CommonLayout'

const statusData = [
  { id: 1, label: 'TODO' },
  { id: 2, label: 'DOING' },
  { id: 3, label: 'DONE' },
  { id: 4, label: 'WAITING' },
]

const TaskItem = ({ item, onPress }) => (
  <TouchableOpacity onPress={onPress} style={styles.item}>
    <Text style={styles.infoText}>
      Name: <Text style={styles.bold}>{item.name}</Text>
    </Text>
    <Text style={styles.infoText}>
      Description: <Text style={styles.bold}>{item.description}</Text>
    </Text>
    <View style={styles.status}>
      {statusData.map((x) => (
        <View
          style={
            x.id === item.status ? styles.activeBadge : styles.inactiveBadge
          }
        >
          <Text
            style={
              x.id === item.status ? styles.activeLabel : styles.inactiveLabel
            }
          >
            {x.label}
          </Text>
        </View>
      ))}
    </View>
  </TouchableOpacity>
)

const ListTask = ({ navigation, onPress }) => {
  const { setLoading, currentProject, loading } = useContext(StoreContext)

  const [data, setData] = useState([])

  const fetchData = useCallback(async () => {
    await setLoading(true)
    const res = await FetchTaskInProject(currentProject.ID)
    await setLoading(false)
    res && setData(res)
  }, [currentProject, setLoading])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  const onAdd = useCallback(() => {
    navigation.navigate('Modal', { screen: 'FormTask' })
  }, [navigation])

  const onDetail = useCallback(() => {
    navigation.navigate('Modal', { screen: 'FormTask' })
  }, [navigation])

  const renderItem = useCallback(
    ({ item, index }) => (
      <TaskItem item={item} onPress={onDetail} key={index.toString()} />
    ),
    [onDetail],
  )
  return (
    <CommonLayout
      data={data}
      renderItem={renderItem}
      filterValue={'name'}
      emptyText={'No task found !'}
      onRefresh={fetchData}
      refreshing={loading}
      onCreate={onAdd}
    />
  )
}

export default ListTask

const styles = StyleSheet.create({
  item: {
    margin: 8,
    borderColor: 'firebrick',
    borderWidth: 2,
    padding: 8,
    borderRadius: 16,
    backgroundColor: 'azure',
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 8,
  },
  bold: {
    fontWeight: 'bold',
  },
  left32: {
    marginLeft: 32,
  },
  top16: {
    marginTop: 8,
  },
  status: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-around',
    paddingVertical: 8,
  },
  activeBadge: {
    flex: 1,
    backgroundColor: 'firebrick',
    marginHorizontal: 4,
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 4,
  },
  inactiveBadge: {
    flex: 1,
    marginHorizontal: 4,
    alignItems: 'center',
    justifyContent: 'center',
    borderColor: 'firebrick',
  },
  activeLabel: {
    color: 'azure',
    fontWeight: 'bold',
  },
  inactiveLabel: {
    fontWeight: 'bold',
  },
  infoText: {
    fontSize: 16,
    margin: 4,
  },
})
