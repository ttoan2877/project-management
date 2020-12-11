import React, { useContext, useState, useCallback, useEffect } from 'react'
import { StoreContext } from '../../App'
import CommonLayout from '../../Components/CommonLayout'
import RoleItem from '../../Components/RoleItem'
import { FetchAllRoleInProject } from '../../Services/RoleService'

const ListRoles = ({ navigation, onPress }) => {
  const { setLoading, loading, currentProject } = useContext(StoreContext)

  const [data, setData] = useState([])

  const fetchData = useCallback(async () => {
    await setLoading(true)
    const res = await FetchAllRoleInProject(currentProject.ID)
    await setLoading(false)
    res && setData(res)
  }, [currentProject, setLoading])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  const onAdd = useCallback(() => {
    navigation.navigate('Modal', { screen: 'CreateRole' })
  }, [navigation])

  const renderItem = useCallback(
    ({ item, index }) => (
      <RoleItem item={item} onPress={onPress} key={index.toString()} />
    ),
    [onPress],
  )

  return (
    <CommonLayout
      data={data}
      renderItem={renderItem}
      filterValue={'name'}
      emptyText={'No role found !'}
      onRefresh={fetchData}
      refreshing={loading}
      onCreate={!onPress && onAdd}
    />
  )
}

export default ListRoles
