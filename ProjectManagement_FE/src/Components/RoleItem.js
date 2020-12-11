import React, { useMemo } from 'react'
import { StyleSheet, Text, View, TouchableOpacity } from 'react-native'

const RoleItem = ({ item, onPress }) => {
  const permission = useMemo(() => {
    let temp = ''
    if (item?.permissions?.includes('1')) {
      return 'Full access'
    }
    if (item?.permissions?.includes('2')) {
      temp += 'Create/Delete task'
    }
    if (item?.permissions?.includes('3')) {
      temp += 'Assign/Unassign member from task'
    }
    if (item?.permissions?.includes('3')) {
      temp += 'Add/Remove member from project'
    }
    return temp
  }, [item?.permissions])

  const Wrapper = onPress ? TouchableOpacity : View

  return (
    <Wrapper onPress={() => onPress({ role_id: item.ID })} style={styles.item}>
      <View style={styles.row}>
        <Text>
          ID: <Text style={styles.bold}>{item.ID}</Text>
        </Text>
        <Text style={styles.left32}>
          Name: <Text style={styles.bold}>{item.name}</Text>
        </Text>
      </View>
      <Text>
        Description: <Text style={styles.bold}>{item.description}</Text>
      </Text>
      <Text style={styles.top16}>
        Permission: <Text style={styles.bold}>{permission}</Text>
      </Text>
    </Wrapper>
  )
}
export default RoleItem

const styles = StyleSheet.create({
  item: {
    margin: 16,
    borderColor: 'firebrick',
    borderWidth: 2,
    padding: 16,
    borderRadius: 16,
    backgroundColor: 'azure',
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 16,
  },
  bold: {
    fontWeight: 'bold',
  },
  left32: {
    marginLeft: 32,
  },
  top16: {
    marginTop: 16,
  },
})
