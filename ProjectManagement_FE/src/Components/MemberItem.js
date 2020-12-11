import React from 'react'
import { StyleSheet, Text, View, TouchableOpacity } from 'react-native'

const MemberItem = ({ data, onPress, onRemove }) => {
  return (
    <TouchableOpacity onPress={onPress} style={styles.item}>
      <View style={styles.info}>
        <Text style={styles.textWrap}>
          ID: <Text style={styles.bold}>{data.ID}</Text>
        </Text>
        <Text style={styles.textWrap}>
          Name: <Text style={styles.bold}>{data.employee.name}</Text>
        </Text>
        <Text style={styles.textWrap}>
          Phone: <Text style={styles.bold}>{data.employee.phone_number}</Text>
        </Text>
        <Text style={styles.textWrap}>
          Email: <Text style={styles.bold}>{data.employee.mail}</Text>
        </Text>
      </View>
      <TouchableOpacity onPress={onRemove} style={styles.btn}>
        <Text style={styles.label}>Remove</Text>
      </TouchableOpacity>
    </TouchableOpacity>
  )
}

export default MemberItem

const styles = StyleSheet.create({
  item: {
    padding: 16,
    borderWidth: 2,
    margin: 8,
    borderColor: 'firebrick',
    backgroundColor: 'azure',
    borderRadius: 8,
    flexDirection: 'row',
    alignItems: 'center',
  },
  info: {
    flex: 1,
  },
  bold: {
    fontWeight: 'bold',
  },
  textWrap: {
    marginTop: 4,
  },
  btn: {
    backgroundColor: 'firebrick',
    paddingVertical: 4,
    paddingHorizontal: 8,
    height: 24,
    margin: 2,
    borderRadius: 8,
  },
  label: {
    color: 'white',
    fontWeight: 'bold',
    fontSize: 12,
  },
})
