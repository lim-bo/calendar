#ifndef CAL_30_H
#define CAL_30_H

#include <QWidget>
#include "client.h"
#include "cfg.h"
#include "event_entry.h"
#include <QHBoxLayout>
#include <QMessageBox>
namespace Ui {
class Cal_30;
}

class Cal_30 : public QWidget
{
    Q_OBJECT

public:
    explicit Cal_30(QString uid, QWidget *parent = nullptr);
    ~Cal_30();

private slots:
    void on_calendarWidget_clicked(const QDate &date);
    void eventDeleted(event_entry*);
private:
    Ui::Cal_30 *ui;
    Client cli;
    QString uid;

};

#endif // CAL_30_H
