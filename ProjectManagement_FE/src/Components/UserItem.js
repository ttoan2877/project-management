import React, { useState, useCallback, Fragment } from 'react'
import { StyleSheet, Text, View, TouchableOpacity } from 'react-native'
import ListRoles from '../Containers/Roles/ListRoles'

const UserItem = ({ data, onPress }) => {
  const [show, setShow] = useState(false)
  const [add, setAdd] = useState(false)
  const onAdd = useCallback(
    ({ role_id }) => {
      return onPress({ role_id, user_id: data.ID })
    },
    [data.ID, onPress],
  )
  return (
    <Fragment>
      <View style={styles.item}>
        <View style={styles.info}>
          <Text>
            ID: <Text style={styles.bold}>{data.ID}</Text>
          </Text>
          <Text>
            Username: <Text style={styles.bold}>{data.username}</Text>
          </Text>
          <Text>
            Name: <Text style={styles.bold}>{data.employee.name}</Text>
          </Text>
          {show && (
            <View style={styles.top4}>
              <Text>
                Bio: <Text style={styles.bold}>{data.employee.bio}</Text>
              </Text>
              <Text>
                Phone:{' '}
                <Text style={styles.bold}>{data.employee.phone_number}</Text>
              </Text>
              <Text>
                Email: <Text style={styles.bold}>{data.employee.mail}</Text>
              </Text>
            </View>
          )}

          {add && (
            <View style={styles.top4}>
              <Text style={styles.bold}>Please select a role below:</Text>
            </View>
          )}
        </View>
        <View>
          <TouchableOpacity onPress={() => setShow(!show)} style={styles.btn}>
            <Text style={styles.label}>{show ? 'Less' : 'More'}</Text>
          </TouchableOpacity>
          <TouchableOpacity onPress={() => setAdd(!add)} style={styles.btn}>
            <Text style={styles.label}>{add ? 'Cancel' : 'Add +'}</Text>
          </TouchableOpacity>
        </View>
      </View>
      {add && <ListRoles onPress={onAdd} />}
    </Fragment>
  )
}

export default UserItem

const styles = StyleSheet.create({
  item: {
    padding: 16,
    borderWidth: 2,
    margin: 8,
    borderColor: 'firebrick',
    backgroundColor: 'azure',
    borderRadius: 8,
    flexDirection: 'row',
  },
  info: {
    flex: 1,
  },
  bold: {
    fontWeight: 'bold',
  },
  btn: {
    backgroundColor: 'firebrick',
    paddingVertical: 4,
    paddingHorizontal: 16,
    height: 24,
    margin: 2,
    borderRadius: 8,
  },
  label: {
    color: 'white',
    fontWeight: 'bold',
    fontSize: 12,
  },
  top4: {
    marginTop: 4,
  },
})
